'use client';

import React, { useState, useEffect } from 'react';
import {
  Card,
  CardBody,
  Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Input,
  Textarea,
  useDisclosure,
  Spinner
} from "@nextui-org/react";

interface Idea {
  id: string;
  title: string;
  description: string;
  created_at: string;
}

export default function IdeaBoardPage() {
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [newIdea, setNewIdea] = useState({ title: '', description: '' });
  const [errors, setErrors] = useState({ title: '', description: '' });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { isOpen, onOpen, onClose } = useDisclosure();

  useEffect(() => {
    fetchIdeas();
  }, []);

  const fetchIdeas = async () => {
    try {
      setIsLoading(true);
      const response = await fetch('/api/ideas');
      if (!response.ok) {
        throw new Error('Failed to fetch ideas');
      }
      const data = await response.json();
      setIdeas(data);
    } catch (error) {
      setError('Failed to load ideas. Please try again later.');
      console.error('Error fetching ideas:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const validateForm = () => {
    let isValid = true;
    const newErrors = { title: '', description: '' };

    if (!newIdea.title.trim()) {
      newErrors.title = 'Title is required';
      isValid = false;
    }

    if (!newIdea.description.trim()) {
      newErrors.description = 'Description is required';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  const handleSubmit = async () => {
    if (validateForm()) {
      try {
        setIsLoading(true);
        const response = await fetch('/api/ideas', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(newIdea),
        });

        if (!response.ok) {
          throw new Error('Failed to create idea');
        }

        const createdIdea = await response.json();
        setIdeas([createdIdea, ...ideas]);
        setNewIdea({ title: '', description: '' });
        onClose();
      } catch (error) {
        setError('Failed to create idea. Please try again.');
        console.error('Error creating idea:', error);
      } finally {
        setIsLoading(false);
      }
    }
  };

  if (isLoading && ideas.length === 0) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold">Idea Board</h1>
        <Button 
          onPress={onOpen}
          color="primary"
          className="px-6"
        >
          Add Idea
        </Button>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {/* Ideas Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {ideas.map((idea) => (
          <Card key={idea.id} className="hover:shadow-lg transition-shadow">
            <CardBody>
              <h3 className="text-lg font-semibold mb-2">{idea.title}</h3>
              <p className="text-gray-600">{idea.description}</p>
              <div className="mt-4 text-sm text-gray-400">
                {new Date(idea.created_at).toLocaleDateString()}
              </div>
            </CardBody>
          </Card>
        ))}
      </div>

      {/* Add Idea Modal */}
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <ModalHeader>Add New Idea</ModalHeader>
          <ModalBody>
            <div className="space-y-4">
              <Input
                label="Title"
                placeholder="Enter idea title"
                value={newIdea.title}
                onChange={(e) => setNewIdea({ ...newIdea, title: e.target.value })}
                isInvalid={!!errors.title}
                errorMessage={errors.title}
              />
              <Textarea
                label="Description"
                placeholder="Describe your idea"
                value={newIdea.description}
                onChange={(e) => setNewIdea({ ...newIdea, description: e.target.value })}
                isInvalid={!!errors.description}
                errorMessage={errors.description}
              />
            </div>
          </ModalBody>
          <ModalFooter>
            <Button color="danger" variant="light" onPress={onClose}>
              Cancel
            </Button>
            <Button 
              color="primary" 
              onPress={handleSubmit}
              isLoading={isLoading}
            >
              Submit
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
}